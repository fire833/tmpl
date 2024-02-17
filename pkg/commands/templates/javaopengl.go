package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewJAVAOPENGLCommand() *cobra.Command {
	type cmdOpts struct {
		Output     *string
		Header     *string
		Package    *string
		Name       *string
		DimensionX *uint
		DimensionY *uint
		FPS        *uint
	}

	const tmpl string = `
{{.Header}}

package {{ .Package }};

import java.awt.Dimension;
import java.awt.Font;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;

import javax.swing.JFrame;
import javax.swing.SwingUtilities;

import com.jogamp.nativewindow.ScalableSurface;
import com.jogamp.opengl.GL;
import com.jogamp.opengl.GL2;
import com.jogamp.opengl.GLAutoDrawable;
import com.jogamp.opengl.GLCapabilities;
import com.jogamp.opengl.GLEventListener;
import com.jogamp.opengl.GLProfile;
import com.jogamp.opengl.awt.GLCanvas;
import com.jogamp.opengl.util.FPSAnimator;
import com.jogamp.opengl.util.awt.TextRenderer;

public final class Application implements GLEventListener, Runnable {
	private int width;
	private int height;

	public static void main(String[] args) {
		SwingUtilities.invokeLater(new Application(args));
	}

	public Application(String[] args) {
	}

	@Override
	public void reshape(GLAutoDrawable drawable, int x, int y, int width, int height) {
		System.out.printf("New window coordinates: (x, y) = (%d, %d), (width, height) = (%d, %d)\n", x, y, width, height);
		this.width = width;
		this.height = height;
	}

	@Override
	public void run() {
		GLProfile profile = GLProfile.getDefault();
		System.out.printf("Java Version %s\nOpenGL Version: %s\nThread number: %d\n\n",
				System.getProperty("java.version"),
				profile.getName(),
				Thread.currentThread().getId());
		// Canvas/frame setup
		GLCapabilities capabilities = new GLCapabilities(profile);
		GLCanvas canvas = new GLCanvas(capabilities); // Single-buffer
		final float[] scale = { ScalableSurface.IDENTITY_PIXELSCALE,
				ScalableSurface.IDENTITY_PIXELSCALE };
		canvas.setSurfaceScale(scale);
		canvas.setPreferredSize(new Dimension({{ .DimensionX }}, {{ .DimensionY }}));

		// JFrame setup
		JFrame frame = new JFrame("{{ .Name }}");
		frame.setBounds(50, 50, 200, 200);
		frame.getContentPane().add(canvas);
		frame.pack();
		frame.setVisible(true);
		frame.setDefaultCloseOperation(JFrame.DISPOSE_ON_CLOSE);
		frame.addWindowListener(new WindowAdapter() {
			// Close out app on window close.
			public void windowClosing(WindowEvent e) {
				System.exit(0);
			}
		});

		canvas.addGLEventListener(this);

		// Running our app at {{ .FPS }} FPS
		FPSAnimator animator = new FPSAnimator(canvas, {{ .FPS }});
		animator.start();
	}

	@Override
	public void init(GLAutoDrawable drawable) {
		this.width = drawable.getSurfaceWidth();
		this.height = drawable.getSurfaceHeight();
		GL2 g = drawable.getGL().getGL2();
		g.glEnable(GL2.GL_POINT_SMOOTH);
	}

	@Override
	public void dispose(GLAutoDrawable drawable) {
	}

	@Override
	public void display(GLAutoDrawable drawable) {
		GL2 g = drawable.getGL().getGL2();
		g.glClear(GL.GL_COLOR_BUFFER_BIT);

		g.glFlush();
	}
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "javaopengl",
		Aliases: []string{"jogl", "javaogl", "jopengl", "javagl"},
		Short:   "Generate boilerplate for creating new JavaOpenGL applications.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("JavaOpenGL").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("JavaOpenGL", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:     set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:     &str,
		Package:    set.StringP("package", "p", "", "Specify the package name for this new application."),
		Name:       set.StringP("name", "n", "Kendall's Homework", "Specify the output name for this JOGL application (what will the window name be?)"),
		DimensionX: set.UintP("dimx", "x", 1280, "Specify the preferred X dimension for created window."),
		DimensionY: set.UintP("dimy", "y", 920, "Specify the preferred Y dimension for created window."),
		FPS:        set.Uint("fps", 60, "Specify the frames per second of your application."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
